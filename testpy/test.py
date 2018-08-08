import requests

url = 'http://ip.chinaz.com/getip.aspx'
kv = {
	"http":"http://89.163.152.44:8080"
}
try:
	r = requests.get(url, proxies = kv)
except:
	print("error")
	exit()
print(r.text)
