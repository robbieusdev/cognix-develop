
### Creating bucket
**when user signin**

The backend system verifies is a bucket exists for current tenant, 
if not it generates a bucket for each tenant in Minio, 
with the bucket name formatted as "tenant-< UUID >"

Files from the public connector will be stored in this bucket directly, 
while files from the private connector will be stored in a folder named "user-< UUID >" within the tenant's bucket.