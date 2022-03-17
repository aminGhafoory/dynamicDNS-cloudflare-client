# dynamicDNS-cloudflare-client

with this client you can make your own private Dynamic DNS server; just like [DuckDNS.org](https://DuckDNS.org) 

this service is crucial if you want to make a homeLab and your ISP does Not provid a static IP address

## **HOW TO USE**
- register a Domain or use and existing one
- create an account on [cloudflare.com](https://cloudflare.com)
- add cloudflare nameservers to your domain at your domain admin panel
- get your apikey and other keys from cloud flare website
- edit `config.toml` file using your own credentials
- compile the program using go `go build -o DDNS`(**Note**: add build tags according to your target os)
  - if you plannig to use this app on linux machines like raspberry pi use this command `GOOS=linux go build -o DDNS`
- add it to cronJobs to run it in 5 minute intervals using a command like this `*/5 * * * * /user/bin/Dyndns DDNS >/dev/null 2>&1` (**Note**: edit file path according to your own system)
- if you setup all thing correctly every time this program runs your domain's Arecord should be updated to your machine's current ip
- have fun 🎈
 
