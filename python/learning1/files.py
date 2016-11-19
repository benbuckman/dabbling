# import pickle
import os


mypath = os.path.realpath(__file__)
print(mypath)


myfile = open(mypath)
# content = pickle.load(myfile, encoding="utf8")
content = myfile.read()
print(content)
myfile.close()
