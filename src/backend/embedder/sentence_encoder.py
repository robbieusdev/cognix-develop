import os
from sentence_transformers import SentenceTransformer
import logging
import threading
from typing import Dict, List
from dotenv import load_dotenv

# compare performance with https://github.com/qdrant/fastembed

# Load environment variables from .env file
load_dotenv()

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class SentenceEncoder:
        # Validate environment variables at class level initialization
    _cache_limit: int = int(os.getenv('MODEL_CACHE_LIMIT', 1))
    _local_model_dir: str = os.getenv('LOCAL_MODEL_PATH', 'models')
    
    if _cache_limit <= 0:
        raise ValueError("MODEL_CACHE_LIMIT must be an integer greater than 0")
    
    # Convert the local model path to an absolute path based on the working directory
    _local_model_dir = os.path.abspath(_local_model_dir)

    if not os.path.isdir(_local_model_dir):
        raise ValueError(f"LOCAL_MODEL_PATH '{_local_model_dir}' is not a valid directory")

    # Thread lock for thread-safe access to the cache
    _lock: threading.Lock = threading.Lock()
    
    # Dictionary to store cached model instances
    _model_cache: Dict[str, SentenceTransformer] = {}


    @classmethod
    def _load_model(cls, model_name: str) -> SentenceTransformer:
        """
        Loads a model from the local directory if available, otherwise downloads and saves it.

        Parameters:
        model_name (str): The name of the model to load or download.

        Returns:
        SentenceTransformer: The loaded SentenceTransformer model.
        """
        model_path: str = os.path.join(cls._local_model_dir, model_name)
        
        if not os.path.exists(model_path) or not os.listdir(model_path):
            logger.info(f"{model_name} model not found locally, downloading from Hugging Face...")
            try:
                model: SentenceTransformer = SentenceTransformer(model_name)
                model.save(model_path)
                logger.info(f"{model_name} model saved locally at {model_path}")
            except Exception as e:
                logger.info(f"❌ {model_name} failed to download or save the model due to: {e}")
        else:
            logger.info(f"loading {model_name} from local directory...")
        
        return SentenceTransformer(model_path)

    @classmethod
    def _get_model(cls, model_name: str) -> SentenceTransformer:
        """
        Retrieves a model from the cache or loads it if not already cached. Manages the cache size.

        Parameters:
        model_name (str): The name of the model to retrieve.

        Returns:
        SentenceTransformer: The model instance.
        """
        with cls._lock:
            # Check if the model is already in the cache
            if model_name in cls._model_cache:
                logger.info(f"using cached model: {model_name}")
                return cls._model_cache[model_name]

            # If the cache limit is reached, unload the oldest model
            if len(cls._model_cache) >= cls._cache_limit:
                oldest_model: str = next(iter(cls._model_cache))
                logger.info(f"unloading model: {oldest_model}")
                # removing model from cache and memory 
                del cls._model_cache[oldest_model]

            # Load and cache the new model
            logger.info(f"loading model: {model_name}")
            model: SentenceTransformer = cls._load_model(model_name)
            cls._model_cache[model_name] = model
            return model

    @classmethod
    def embed(cls, text: str, model_name: str) -> List[float]:
        """
        Encodes the provided text using the specified model.

        Parameters:
        text (str): The text to be encoded.
        model_name (str): The name of the model to use for encoding.

        Returns:
        list: A list of floats representing the encoded text.
        """
        model: SentenceTransformer = cls._get_model(model_name)
        return model.encode(text).tolist()