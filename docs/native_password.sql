ALTER USER 'apps'@'%' IDENTIFIED WITH caching_sha2_password BY 'apps';
GRANT ALL PRIVILEGES ON kf_ai.* TO 'apps'@'%' IDENTIFIED BY 'apps';
FLUSH PRIVILEGES;