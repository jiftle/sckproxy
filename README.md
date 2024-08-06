# sckproxy

simple network proxy

## proto

- socket5
- http


## useage

```
# socket5 proxy
sckproxy -l :1080
sckproxy -l :1080 -m socket

# http proxy
sckproxy -l :1080 -m http
```