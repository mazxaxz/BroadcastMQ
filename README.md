[![Contributors][contributors-shield]][contributors-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![Apache-2.0 License][license-shield]][license-url]



<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/mazxaxz/BroadcastMQ">
    <img src="images/bmq-logo.png" alt="Logo" height="180">
  </a>

  <h3 align="center">BroadcastMQ</h3>

  <p align="center">
    A lightweight app for Broadcasting events across your Brokers!
    <br />
    <br />
    <a href="https://hub.docker.com/r/mazxaxz/broadcast-mq">Dockerhub</a>
    ·
    <a href="https://github.com/mazxaxz/BroadcastMQ/issues">Report a bug</a>
    ·
    <a href="https://github.com/mazxaxz/BroadcastMQ/issues">Request feature</a>
  </p>
</p>



<!-- TABLE OF CONTENTS -->
## Table of Contents

* [About the Project](#about-the-project)
  * [Built With](#built-with)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
  * [Configuration](#configuration)
* [Usage](#usage)
* [Roadmap](#roadmap)
* [Contributing](#contributing)
* [License](#license)
* [Contact](#contact)
* [Acknowledgements](#acknowledgements)



<!-- ABOUT THE PROJECT -->
## About The Project

It is a simple lightweight application for broadcasting messages across multiple message brokers based on your configuration. It can be hosted anywhere since it is contenerized. Just configure and you are good to go.

Why did I created this?
Many of you might have experienced a situation, where your application consumes tons of messages every second on the production environment, right?  
How about other environments? Is your application exhausted with amount of events or maybe it's bored as hell?

I have created BroadcastMQ to solve the lack of events on the than production environments.  
Why there are multiple development environments, when there is no actual data. Everything we build is data oriented.  
Of course, if your application consumes **sensitive** data I **would not** recommend using it, because of the security things etc.

### Built With
* [Golang](https://golang.org)
* [Docker](https://www.docker.com)
* [Streadway/amqp](https://github.com/streadway/amqp)
* [Sirupsen/logrus](https://github.com/sirupsen/logrus)
* [Kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig)
* [Guiferpa/gody](https://github.com/guiferpa/gody)
* [Segmentio/ksuid](https://github.com/segmentio/ksuid)



<!-- GETTING STARTED -->
## Getting Started

If you want to run BroadcastMQ locally, **Docker** is a must.  
Having docker installed, you can just: `docker pull mazxaxz/broadcast-mq:latest`

### Prerequisites

* docker
[Get docker](https://docs.docker.com/get-docker/)

### Installation

1. `git clone https://github.com/mazxaxz/BroadcastMQ.git`
2. `docker-compose -f ./examples/docker-compose.yaml up --build`


### Configuration

* environment variables
```
BMQ_CONFIGPATH=<path>                                     default:"/etc/broadcastmq/config.yaml"
BMQ_LOGLEVEL=(trace|debug|info|warn|error|fatal|panic)    default:"info"
BMQ_OUTPUTTYPE=(text|json)                                default:"text"
```

* config file
Configuration possibilities are provided in `examples/config.yaml` file [here](https://github.com/mazxaxz/BroadcastMQ/blob/master/examples/config.yaml)


<!-- USAGE EXAMPLES -->
## Usage

Example using docker-compose can be seen [here](https://github.com/mazxaxz/BroadcastMQ/blob/master/examples)

Steps:
1. `git clone https://github.com/mazxaxz/BroadcastMQ.git`
2. `docker-compose -f ./examples/docker-compose.yaml up --build`
3. Go to `localhost:25673` and publish message onto `MQ.Topic.Source.Exchange` exchange with `MQ.RoutingKey` routing key
4. Go to `localhost:35673`
5. If `MQ.Queue.Destination.Example.1` queue contains the message you have published earlier, it means BroadcastMQ is working


<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/mazxaxz/BroadcastMQ/issues) for a list of proposed features (and known issues).



<!-- CONTRIBUTING -->
## Contributing

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the Apache-2.0 License. See `LICENSE` for more information.



<!-- CONTACT -->
## Contact

Adrian Ciura - [@mazxaxz](https://twitter.com/mazxaxz)

Project Link: [github.com/mazxaxz/BroadcastMQ](https://github.com/mazxaxz/BroadcastMQ)

Dockerhub: [hub.docker.com/r/mazxaxz/broadcast-mq](https://hub.docker.com/r/mazxaxz/broadcast-mq)



<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements
* [Freepik](http://www.freepik.com)
* [Flaticon](https://www.flaticon.com)
* [Img Shields](https://shields.io)
* [othneildrew](https://github.com/othneildrew/Best-README-Template)
* [Norbert Włodarczyk](https://github.com/RvuvuzelaM)





<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/mazxaxz/BroadcastMQ.svg?style=flat-square
[contributors-url]: https://github.com/mazxaxz/BroadcastMQ/graphs/contributors
[stars-shield]: https://img.shields.io/github/stars/mazxaxz/BroadcastMQ.svg?style=flat-square
[stars-url]: https://github.com/mazxaxz/BroadcastMQ/stargazers
[issues-shield]: https://img.shields.io/github/issues/mazxaxz/BroadcastMQ.svg?style=flat-square
[issues-url]: https://github.com/mazxaxz/BroadcastMQ/issues
[license-shield]: https://img.shields.io/github/license/mazxaxz/BroadcastMQ.svg?style=flat-square
[license-url]: https://github.com/mazxaxz/BroadcastMQ/blob/master/LICENSE.txt