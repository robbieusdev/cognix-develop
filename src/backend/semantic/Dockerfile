# # ###################################
# # Distroless working
# # ###################################
# FROM python:3.11.7-slim  AS build-env

# WORKDIR /app
# ADD ./ ./ 

# # Install any needed packages specified in requirements.txt
# RUN pip install --upgrade pip && \
#     pip install --no-cache-dir -r requirements.txt \
#     && rm -rf /root/.cache

# # Install dependencies and clean up
# RUN apt-get update && apt-get install -y --no-install-recommends \
#     && rm -rf /var/lib/apt/lists/*

# # Use a distroless Python image as a parent image
# FROM gcr.io/distroless/python3-debian12

# # Copy the application code and installed packages from the build environment
# COPY --from=build-env /app /app
# COPY --from=build-env /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages

# # Set the working directory in the container
# WORKDIR /app

# # Set environment variable for Python path
# ENV PYTHONPATH=/usr/local/lib/python3.11/site-package

# # Run the application
# CMD ["semantic_service.py"]





# ###################################
# python slim working
# ###################################
# using Python 3.11.0 as the newr supported by PyTorch
FROM python:3.11.7-slim-bookworm


# Set the working directory inside the container
WORKDIR /app

# Copy the requirements.txt first to leverage Docker cache
COPY requirements.txt .

# Install dependencies and clean up
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*
RUN pip install --upgrade pip
# RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
RUN pip install --no-cache-dir -r requirements.txt \
    && rm -rf /root/.cache

# Copy the rest of the embedder application code
COPY . .

# Command to run your application
CMD ["python", "semantic_service.py"]













# # Use a slim Python base image
# # using Python 3.11.0 as the newr supported by PyTorch
# FROM python:3.11.7-slim-bookworm


# # Set the working directory inside the container
# WORKDIR /app

# # Copy the requirements.txt first to leverage Docker cache
# COPY requirements.txt .

# # Install dependencies and clean up
# RUN apt-get update && apt-get install -y --no-install-recommends \
#     build-essential \
#     && rm -rf /var/lib/apt/lists/*
# RUN pip install --upgrade pip
# RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
# RUN pip install --no-cache-dir -r requirements.txt \
#     && rm -rf /root/.cache

# # Copy the rest of the embedder application code
# COPY . .

# # Command to run your application
# CMD ["python", "semantic.py"]
