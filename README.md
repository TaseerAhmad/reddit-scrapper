# reddit-scrapper

Mini scrapper. Useful for scrapping on the Go.





### Installation



    $ git clone https://github.com/TaseerAhmad/reddit-scrapper.git

    $ cd reddit-scrapper

    $ go run main.go arg[1] arg[2] arg[3]

    

     

### Argument Defination


- `arg[1]` Evaluates to `true` or `false`. When passed `"true"`, the scrapper process will not terminate

 

- `arg[2]` Evaluates to an integer used for defining an interval at which the scrapper restarts. Defined in seconds.



- `arg[3]` Evaluates the URL used for scrapping. The URL passed must be the subreddit's URL.





### Example

    go run main.go true 5 https://old.reddit.com/r/worldnews/new //Will not terminate and keep scrapping after every 5 seconds.

    go run main.go false 5 https://old.reddit.com/r/worldnews/new //Will keep scrap only 1 time and terminate.

    

    



### TODO

- Use `termui` for the continuous scrapping function

- Use `Cobra` for advanced CLI interaction

- Introduce a tag watching feature

- Implement a bot to send Telegram messages
