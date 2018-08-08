import requests

local = "http://127.0.0.1:2333/proxy?act=get"
r = requests.get(local)
if r.json()["code"] == 200:
    ip = r.json()["proxies"][0]["ip"]
    port = r.json()["proxies"][0]["port"]

url = "http://nemesisly.xyz"
proxy = {
    "http":"http://" + ip + ":" + port
}
res = requests.get(url, proxies = proxy, timeout = 5)
print(res.text)
