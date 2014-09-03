books
=====

command-line book inventory management

**WARNING**: this software is in development. Don't use it unless you know what you're doing. Which is also why I don't write installation instructions, yet.

Database
--------

When adding your first entry, a Sqlite database will be created in *~/.books/books/db*

Example usage
-------------

    $ books add
    Author: Frank Herbert
    Title: Dune
    Comments:

    $ books add
    Author: Stephen King
    Title: The Stand
    Comments:
    
    $ books ls
    {#1: 'Dune' by 'Frank Herbert', comments: }
    {#2: 'The Stand' by 'Stephen King', comments: }

    $ books ls sta
    {#2: 'The Stand' by 'Stephen King', comments: }
    
    $ books ls herb
    {#1: 'Dune' by 'Frank Herbert', comments: }
    
    $ books help
    USAGE: books COMMAND argument1 argument2 ...
    Available commands:
        	ls - list books, pass search terms as arguments
        	add - add books
        	help - display this text
    

Future ideas
------------

* ~~delete books~~
* edit books
* use the [OpenLibrary](https://openlibrary.org/developers/api) and/or [Librarything](https://www.librarything.com/services/) and/or [Goodreads](https://www.goodreads.com/api) API to fill in ISBN, maybe autocomplete/autofix titles/authors etc.
* implement "lending" (to whom, set a deadline, show reminders)
* web gui
* export to csv/tsv. What's the industry standard format?
