# freedns-go

Fork from:https://github.com/Chenyao2333/freedns-go

Add read china ip list file instead of integrate ip list to program.

## Usage

```
sudo ./freedns-go -f 114.114.114.114:53 -c 8.8.8.8:53 -l 0.0.0.0:53 -r /root/chnroute.txt
```
Important: Don't forget port parameter in every dns address, whatever is 53 default
