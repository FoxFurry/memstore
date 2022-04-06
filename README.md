<div id="top"></div>


<!-- PROJECT SHIELDS -->

[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]


<!-- PROJECT LOGO -->
<br />
<div align="center">
  
<h3 align="center">memstore</h3>

  <p align="center">
    A simple key-value in-memory store with built-in sharding and load-balancing
    <br /> 
  </p>
</div>



<!-- ABOUT THE PROJECT -->
## About The Project

**This project is my personal project, and it does not offer security, consistency and other production essential feature. Do not use it in live environment**

memstore is a simple key-value in-memory storage which uses REST API as main interface. The minimum operator of memstore is transaction
and transactions are atomic and isolated (like SQL transaction). Every transaction is executed multithreaded on a copy of a shard and after successful execution of all commands -
transaction is being added to the shard queue for actual apply. Right now supported operations are only `GET` and `SET`. Despite little amount of operations, store does automatic shardering and load-balancing between
shards. Default number of shards if 4, but in future I will add mechanism to increase or decrease amount of shards on fly. Every operation inside of transaction is mapped to
exact shard using consistent hashing. The project is in working condition, but still a lot of work is required

<p align="right">(<a href="#top">back to top</a>)</p>



## Dependencies

* [github.com/buraksezer/consistent](https://github.com/buraksezer/consistent)
* [github.com/cespare/xxhash](https://github.com/cespare/xxhash)
* [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* [github.com/google/btree](https://github.com/google/btree)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.



### Prerequisites

* Install Go 1.18 (or at least 1.15)  
  [Go Install Guide](https://golang.org/doc/install)

  

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/FoxFurry/memstore
   ```
2. Install the dependencies
   ```shell
    $ go mod download
    ```

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage

To launch the project after copying repo, use next
```shell
$ go run main.go serve --port=8000
```
Port flag can be omitted, default value is 8080

Documentation: https://pkg.go.dev/github.com/FoxFurry/memstore

<p align="right">(<a href="#top">back to top</a>)</p>




<!-- ROADMAP -->
## Roadmap

- [x] Basic command execution
- [x] Node implementation
  - [x] Node itself
  - [x] Node journal mechanism
  - [x] Node copying mechanism
- [x] Cluster implementation
  - [x] Consistent hashing
- [x] Basic web interface
- [x] Add cobra CLI interface
- [ ] Add config files
- [ ] Fix SET/GET data race
- [ ] Fix strange `btree.Item is nil` issue
- [ ] Improve performance
  - [ ] Under 100ns for GET
  - [ ] Under 200ns for SET
- [ ] To be continued ...

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>







<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [Table and Index Organization](http://msdn.microsoft.com/en-us/library/ms189051.aspx)
* [Planning and Architecture (Database Engine)](https://docs.microsoft.com/en-us/previous-versions/sql/sql-server-2008-r2/cc280361(v=sql.105))

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
[license-shield]: https://img.shields.io/github/license/FoxFurry/memstore.svg?style=for-the-badge
[license-url]: https://github.com/FoxFurry/memstore/blob/master/LICENSE

[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/arthur-isac-412a6519b/

[for-love-shield]: https://forthebadge.com/images/badges/built-with-love.svg
[for-love-url]: https://github.com/FoxFurry

[optimized-shield]: https://forthebadge.com/images/badges/0-percent-optimized.svg
[optimized-url]: https://docs.microsoft.com/en-us/previous-versions/sql/sql-server-2008-r2/cc280361(v=sql.105)