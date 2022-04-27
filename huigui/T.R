data = read.csv("./PlantM.txt")

res = lm( data$X1~data$X0 )

print(data$X1)
print(data$X0)
print(res)

print( predict(res, data.frame(X0=12)) )
