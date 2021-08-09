# reddit-scrapper

Mini scrapper. Useful for scrapping on the Go.





### Installation



    $ git clone https://github.com/TaseerAhmad/reddit-scrapper.git

    $ cd reddit-scrapper

    $ go run main.go arg[1] arg[2] arg[3]



### Argument Defination


- `arg[1] int` The number of pages to crawl



- `arg[2] string` The subreddit to scrap



- `arg[3] string` The name of the file to be used in storing the scrapped data



### Example

    go run main.go 2 https://old.reddit.com/r/worldnews/new worldnews.json //Produces a JSON containing all the scrapped data