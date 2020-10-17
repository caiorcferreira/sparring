# Sparring

<p style="text-align: center">The best partner for hammering your application during load testing.</p>

## Usage
Sparring only needs its configuration file to run. You can define the file location with the `SPARRING_CONFIG` environment variable.
  
For example, the following file will create two targets and tells Sparring to run at port `9001` with the `warn` log level.

```yaml
port: 9001
logLevel: warn

targets:
  - method: GET
    path: "/hero"
    statusCode: 200
    body:
      file: samples/get-hero.json
    responseTime: 1s
  - method: GET
    statusCode: 404
    path: "/sidekick"
    responseTime: 1s
```

## Contributing

Please feel free to open issues and pull requests! 

## License

Distributed under the MIT License.
