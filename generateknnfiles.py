import requests as req
# f=open('./docs/listar.json','w',encoding='utf-8')
# f.write(req.get('http://localhost/listar').text)
# f.close()
for i in range(72,101):
    f=open('./docs/responseK'+str(i)+'.json','w',encoding='utf-8')
    f.write(req.get('http://localhost/knn?k='+str(i)).text)
    f.close()