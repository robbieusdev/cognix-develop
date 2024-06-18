import asyncio
import logging

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

async def main():
    try:
        while True:
            await asyncio.sleep(1)
            logger.info("embedder is working")
    except KeyboardInterrupt:
        logger.info("keyb interrupt")

if __name__ == "__main__":
    asyncio.run(main())