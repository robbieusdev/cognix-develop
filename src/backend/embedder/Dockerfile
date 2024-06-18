# ###################################
# Distroless working
# ###################################
FROM python:3.11.7-slim  AS build-env

WORKDIR /app
ADD ./ ./ 

RUN pip3 install --upgrade pip && \
    pip install -r ./requirements.txt

# Install any needed packages specified in requirements.txt
RUN pip install --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt \
    && rm -rf /root/.cache

# Install dependencies and clean up
RUN apt-get update && apt-get install -y --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

# the container before this 3 lines is 107mb, after it's 982mb!!!!
RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
RUN pip install transformers==4.40.2 tqdm==4.66.4 numpy==1.26.4 scikit-learn==1.4.2 scipy==1.13.0 nltk==3.8.1 pillow==10.3.0
RUN pip install --no-deps sentence-transformers==2.7.0

FROM gcr.io/distroless/python3-debian12
COPY --from=build-env /app /app
COPY --from=build-env /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages


WORKDIR /app
ENV PYTHONPATH=/usr/local/lib/python3.11/site-packages
CMD ["embedder_service.py"]







# # ###################################
# # python slim working
# # ###################################
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
# CMD ["python", "embedder_service.py"]














# # Use a multi-stage build to ensure the final image has the correct architecture
# # Stage 1: Build stage
# FROM --platform=linux/amd64 python:3.11-slim as build-env

# # Set the working directory inside the container
# WORKDIR /app

# USER root
# # Install necessary system packages for grpcio
# RUN apt-get update && apt-get install -y --no-install-recommends \
#     build-essential \
#     libssl-dev \
#     libffi-dev \
#     python3-dev \
#     python3.11-distutils \
#     zlib1g \
#     zlib1g-dev \
#     libexpat1-dev \
#     && rm -rf /var/lib/apt/lists/*

# # Copy the requirements.txt first to leverage Docker cache
# COPY requirements.txt .

# # Install dependencies and clean up
# RUN pip install --upgrade pip
# RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
# RUN pip install --no-cache-dir --use-deprecated=legacy-resolver  -r requirements.txt \
#     && rm -rf /root/.cache

# # Copy the rest of the embedder application code
# COPY . /app

# # Stage 2: Production stage
# FROM --platform=linux/amd64 gcr.io/distroless/python3-debian12

# # Set the working directory inside the container
# WORKDIR /app

# # Copy the application from the build stage
# COPY --from=build-env /app /app


# # Required by lots of packages - e.g. six, numpy, wsgi
# COPY --from=build-env /lib/x86_64-linux-gnu/libz.so.1 /lib/x86_64-linux-gnu/
# # Required by google-cloud/grpcio
# COPY --from=build-env /usr/lib/x86_64-linux-gnu/libffi* /usr/lib/x86_64-linux-gnu/
# COPY --from=build-env /lib/x86_64-linux-gnu/libexpat* /lib/x86_64-linux-gnu/
# COPY --from=build-env /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.9/site-packages/
# COPY --from=build-env /usr/local/lib/python3.11/site-packages/grpc /usr/local/lib/python3.9/site-packages/grpc/


# # Set the environment variable for Python path
# ENV PYTHONPATH=/app

# # Command to run your application
# CMD ["embedder.py"]












# # # Use a slim Python base image
# # # using Python 3.11.0 as the newr supported by PyTorch
# # FROM python:3.11.7-slim-bookworm as build-env


# # # Set the working directory inside the container
# # WORKDIR /app

# # # Copy the requirements.txt first to leverage Docker cache
# # COPY requirements.txt .

# # # Install dependencies and clean up
# # RUN apt-get update && apt-get install -y --no-install-recommends \
# #     build-essential \
# #     && rm -rf /var/lib/apt/lists/*
# # RUN pip install --upgrade pip
# # RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
# # RUN pip install --no-cache-dir -r requirements.txt \
# #     && rm -rf /root/.cache

# # # Copy the rest of the embedder application code
# # COPY . .

# # # Command to run your application

# # # Stage 2: Production stage
# # FROM --platform=linux/amd64 gcr.io/distroless/python3-debian12

# # COPY --from=build-env /app /app
# # # Required by lots of packages - e.g. six, numpy, wsgi
# # COPY --from=build-env /lib/x86_64-linux-gnu/libz.so.1 /lib/x86_64-linux-gnu/
# # # Required by google-cloud/grpcio
# # COPY --from=build-env /usr/lib/x86_64-linux-gnu/libffi* /usr/lib/x86_64-linux-gnu/
# # COPY --from=build-env /lib/x86_64-linux-gnu/libexpat* /lib/x86_64-linux-gnu/

# # #COPY --from=build-env /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
# # WORKDIR /app
# # ENV PYTHONPATH=/app
# # CMD ["embedder.py"]








# # crazy
# # https://alex-moss.medium.com/creating-an-up-to-date-python-distroless-container-image-e3da728d7a80
# #