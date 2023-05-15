# Wind Guide
A simple service registry and discovery center written by ARPC.

## API Documentaion

### Register Service
- #### Endpoint
    /register-service

- #### Request Body
  ```
    message RegisterRequest {
        string service_id = 1;
        string service_name = 2;
        string service_addr = 3;
        string service_port = 4;
        string service_version = 5;
        string unique_id = 6;
        string health_check_url = 7;
        int64 usage_count = 8;
    }
  ```

- #### Response Parameters
  ```
    message RegisterResponse {
        string code = 1;
        string message = 2;
    }
  ```

### Discovery Service
- #### Endpoint
    /discovery-service

- #### Request Body
  ```
    message DiscoveryRequest {
        string service_name = 1;
        string version = 2;
        string unique_id = 3;
        string caller_service_name = 4;
        string caller_service_version = 5;
        string caller_unique_id = 6;
        string caller_service_addr = 7;
        string caller_service_port = 8;
    }
  ```

- #### Response Parameters
  ```
    message RegisterRequest {
        string service_id = 1;
        string service_name = 2;
        string service_addr = 3;
        string service_port = 4;
        string service_version = 5;
        string unique_id = 6;
        string health_check_url = 7;
        int64 usage_count = 8;
    }
  ```