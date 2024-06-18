from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from webdriver_manager.chrome import ChromeDriverManager
from bs4 import BeautifulSoup
import time
import logging
from lib.spider.chunked_item import ChunkedItem
from typing import List
from urllib.parse import urljoin, urlparse

from readiness_probe import ReadinessProbe


class SeleniumSpider:

    def __init__(self, base_url):
        self.visited = set()
        self.collected_data: List[ChunkedItem] = [] 
        self.base_domain = urlparse(base_url).netloc
        self.logger = logging.getLogger(__name__)

        # Initialize WebDriver
        chrome_options = Options()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        self.driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()), options=chrome_options)

    def fetch_and_parse(self, url):
        try:
            self.logger.info(f"Processing URL: {url}")
            # Open the URL
            self.driver.get(url)
            # Wait for the dynamic content to load
            time.sleep(10)
            # Get the page source after all scripts have been executed
            html = self.driver.page_source
            # Parse the page with BeautifulSoup
            soup = BeautifulSoup(html, 'html.parser')
            return soup
        except Exception as e:
            self.logger.error(f"❌ Error fetching URL: {url}, Error: {e}")
            return None

    def extract_data(self, soup):
        elements = soup.find_all(['p', 'article', 'div'])
        paragraphs = []

        for element in elements:
            text = element.get_text(strip=True)
            if text:
                paragraphs.append(text)

        formatted_text = '\n\n'.join(paragraphs)
        return formatted_text

    def process_page(self, url: str, recursive: bool) -> list[ChunkedItem] | None:
        start_time = time.time()

        # notifying the readiness probe that the service is alive
        ReadinessProbe().update_last_seen()

        if url in self.visited:
            return

        self.visited.add(url)
        soup = self.fetch_and_parse(url)
        if not soup:
            return

        page_content = self.extract_data(soup)
        if page_content:
            self.collected_data.append(ChunkedItem(url=url, content=page_content))

        links = [a['href'] for a in soup.find_all('a', href=True)]
        for link in links:
            absolute_link = urljoin(url, link)
            parsed_link = urlparse(absolute_link)
            if parsed_link.scheme in ['http', 'https'] and absolute_link not in self.visited:
                if parsed_link.netloc == self.base_domain:
                    self.process_page(absolute_link)

        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        self.logger.info(f"Total elapsed time: {elapsed_time:.2f} seconds")

        # Return the collected data only after all recursive calls are complete
        return self.collected_data

    def close(self):
        self.driver.quit()

# if __name__ == "__main__":
#     logging.basicConfig(level=logging.INFO)
#     spider = SeleniumSpider("https://example.com")
#     spider.process_page("https://example.com")
#     collected_data = spider.get_collected_data()
#     for data in collected_data:
#         print(data.url)
#         print(data.content)
#     spider.close()

# from selenium import webdriver
# from selenium.webdriver.chrome.service import Service
# from selenium.webdriver.chrome.options import Options
# from webdriver_manager.chrome import ChromeDriverManager
# from bs4 import BeautifulSoup
# import time
# import logging
# from lib.chunked_list import ChunkedList
# from typing import List
# from urllib.parse import urljoin, urlparse
# # pip install selenium
# # pip install webdriver-manager
# # pip install beautifulsoup4


# class SeleniumSpider:

#     def __init__(self, base_url):
#         self.visited = set()
#         self.collected_data: List[ChunkedList] = [] 
#         self.base_domain = urlparse(base_url).netloc
#         self.logger = logging.getLogger(__name__)

#     def get_collected_data(self):
#         return self.collected_data

#     def process_page(self, url):
#         try:
#             self.logger.info("creating webdriver..")
        
#             # Setup Chrome WebDriver in headless mode
#             self.logger.info("Setting up browser options")
#             chrome_options = Options()
#             chrome_options.add_argument("--headless")  # Ensures Chrome runs in headless mode
#             chrome_options.add_argument("--no-sandbox")  # Bypass OS security model, mandatory on some systems
#             chrome_options.add_argument("--disable-dev-shm-usage")  # Overcome limited resource problems
#             driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()), options=chrome_options)

#             try:
#                 # Open the URL
#                 driver.get(url)
                
#                 # Wait for the dynamic content to load, adjust the wait time as necessary
#                 time.sleep(10)  # Consider using WebDriverWait for a more efficient wait
                
#                 # Get the page source after all scripts have been executed
#                 html = driver.page_source
                
#                 # Parse the page with BeautifulSoup
#                 soup = BeautifulSoup(html, 'html.parser')
                
#                 # Extract and print text from each paragraph
#                 paragraphs = soup.find_all('p')
#                 for paragraph in paragraphs:
#                     self.logger.info(paragraph.text)
#             finally:
#                 # Make sure to close the driver
#                 driver.quit()
#         except Exception as e:
#             self.logger.exception(e)
        
