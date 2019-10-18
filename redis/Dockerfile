FROM redis:latest

# set cache to 512M and set policy to lru
RUN echo maxmemory 512mb >> /etc/redis.conf && echo maxmemory-policy allkeys-lru >> /etc/redis.conf
CMD ["redis-server", "/etc/redis.conf"]
