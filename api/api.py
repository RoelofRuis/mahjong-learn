from urllib import request, parse

data = {"1": 1}
data = parse.urlencode(data).encode()

req = request.Request("http://localhost:8000/game/1", method="POST", data=data)
resp = request.urlopen(req)

print(resp)
