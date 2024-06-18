import os

# todo: change local path to S3 storage
class SentenceEncoder:
    def __init__(self, model_name, local_model_dir='models'):
        """
        Initializes an instance of SentenceEncoder, attempting to load the model from a local directory first.
        If the model is not available locally, it downloads from Hugging Face and saves it locally.

        Parameters:
        model_name (str): The name of the model to load or download.
        local_model_dir (str): The directory to check for the model and to save the model.
        """
        self.model_path = os.path.join(local_model_dir, model_name)
        


    def embed(self, text):
        """
        Encodes the provided text using the loaded SentenceTransformer model.
        
        Parameters:
        text (str): The text to be encoded.
        
        Returns:
        list: A list of floats representing the encoded text.
        """
        # Use the loaded model to encode the text
        return ""

# # Example usage
# if __name__ == "__main__":
#     model_name = 'sentence-transformers/paraphrase-multilingual-mpnet-base-v2'
#     encoder = TextEncoder(model_name)  # Create an instance of TextEncoder with a specific model
#     encoded_data = encoder.embed("explain routed events in WPF")  # Call the embed method with a sample text
#     print(encoded_data)  # Print the encoded data
