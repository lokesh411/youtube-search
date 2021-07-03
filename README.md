## Installation
1. Clone this repository
2. Create a Key in google API console [GCP Console](https://console.cloud.google.com/apis/dashboard) and add that key in youtube_access_tokens (seperated by comma's if you have more than 1) env variable in docker-compose.yaml file
3. Make sure that docker and docker-compose is installed and running
4. Run docker compose
```bash
  docker-compose up
```
## Design Process
1. One api Key has 10k units of threshold, Search api actually takes 100 units, So, we can store the values in redis and then if the quota of one api key exceeds then use the other, use pacific time to calculate the expiry
2. An index in published_time which will facilitate sorting
3. Full text index on title to allow full text search
## API Endpoints
1. To Search for items
```bash
  curl --location --request GET 'localhost:9000/videos/search?term=arjuna%20triple'
```
2. To Fetch all the items
```bash
  curl --location --request GET 'localhost:5000/videos?publishedDate=2021-07-02T14:36:14Z' // Fetches all the videos in the decending order of published date lesser than the given one, If publishedDate is not present, It fetches the first 10 items with decending order of publishedDate
```
## Uninstallation
```bash
  docker-compose down
```

## MYSQL Table structure
```
  CREATE TABLE `videos` (\n  `id` bigint unsigned NOT NULL AUTO_INCREMENT,\n  `created_at` datetime(3) DEFAULT NULL,\n  `updated_at` datetime(3) DEFAULT NULL,\n  `deleted_at` datetime(3) DEFAULT NULL,\n  `title` varchar(191) DEFAULT NULL,\n  `description` longtext,\n  `published_time` datetime(3) DEFAULT NULL,\n  `thumbnails` json DEFAULT NULL,\n  PRIMARY KEY (`id`),\n  KEY `idx_videos_deleted_at` (`deleted_at`),\n  KEY `idx_videos_published_time` (`published_time`),\n  FULLTEXT KEY `idx_videos_title` (`title`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

## Things that can be improved
1. On a large scale Elastic search can be used for full text search on title and description
2. Dashboard can be built to seemlessly show data
3. Implementing pagination on search apis
