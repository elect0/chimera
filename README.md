# chimera

chimera is a high concurrency service written in go for dealing with on-the-fly image processing.

## installation

1.  **clone**:
    ```bash
    git clone https://github.com/elect0/chimera.git
    ```
2.  **configration**:
    ```bash
    mkdir -p configs && touch configs/config.yaml
    ```
    *config example:*
    ```yaml
    http_server:
      port: 8080
      shutdown_timeout: 5s
    
    log:
      level: "info"
    
    s3:
      # fill in your details about your bucket
      bucket: "your_bucket"
      region: "your_region"
    
    redis:
      address: "localhost:6379"
      password: ""
      db: 0
    
    security:
      # generate a strong key
      hmac_enabled: false
      hmac_secret_key: "secret_key"
      remote_fetch:
        max_download_size_mb: 25
    ```
3.  **run the stack**
    1. **start the dependencies**
      ```bash
      docker-compose up -d
      ```
    2. **start the chimera app**
      ```bash
      go run ./cmd/chimera
      ```
      you should now have the following services running:
      chimera: `http://localhost:8080`
      prometheus: `http://localhost:9091`
      grafana: `http://localhost:3000`
      
5. **try a request**
   all `/transform` requests have to be signed. for local testing, temporarily disable this by setting the        `hmac_enabled` field to false in your `config.yaml`
   example:
   fetching from the s3 bucket:
   ```bash
   curl -o testing.jpg "https://localhost:8080/transform?path=your-image-in-s3-bucket.jpg&width=500"
   ```
   or from a public url:
   ```bash
   curl -o remote-test.jpg "https://localhost:8080transform?url=https%3A%2F%2Fimages.unsplash.com%2Fphoto-1543852786-1cf6624b9987&width=500"
   ```

## api reference

`GET /transform`
| parameter | type | required | description | example |
|---|---|---|---|---|
| `path` | string | **Yes** (or `url`) | object key of the image in the s3 bucket | `my-folder/image.jpg` |
| `url` | string | **Yes** (or `path`) | public url to an image | `https%3A%2F%2F...` |
| `width` | int | No | the target width in pixels | `500` |
| `height`| int | No | the target height in pixels | `300` |
| `quality`| int | No | the quality of the output (1-100) | `85` |
| `crop` | string | No | the crop strategy. use `smart` for saliency-based cropping. | `smart` |
| `watermark`| string | No | the path to a watermark image in your s3 bucket | `logo.png` |
| `wm_pos`| string | No | position of the watermark | `south-east` |
| `wm_opacity`| float | No | opacity of the watermark (0.0-1.0) | `0.7` |
| `s` | string | **Yes** (if enabled) | HMAC-SHA256 signature of the request | `a1b2c3...` |

## roadmap

the project is still underdeveloped. the next major things are:
- complete observability: grafana dashboards
- full cloud-native deployment

## license

[MIT](https://choosealicense.com/licenses/mit/)
  
   
    
