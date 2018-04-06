# go-fibo

This is a go-fibo microservice which just generates load on CPU in the container

To deploy this application, execute the following commands:

  1. Clone repo:

```
$ git clone https://github.com/damianjaniszewski/go-fibo
$ cd go-fibo
```

  2. Create Docker image

```
$ docker build -t <repo>/<name>:<version> .
```

  3. Push Docker image to repository

```
$ docker push <repo>/<name>:<version>
```

  4. Deploy to K8s (verify which image you use)

```
$ kubectl apply -f go-fibo.yml
```

  5. Get go-fibo service port

```
$ kubectl get svc go-fibo
```

  6. Generate load with Apache Bench

```
$ ab -n 2048 -c 128 -s 3600 -m POST http://<addr>:<port>/42
```
