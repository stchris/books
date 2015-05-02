books
=====

command-line book inventory management

**WARNING**: this software is in development. Don't use it unless you know what you're doing. Which is also why I don't write installation instructions, yet.

[![Build Status](https://travis-ci.org/stchris/books.png?branch=master)](https://travis-ci.org/stchris/books)

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

    $ books web
    2015/05/02 13:32:49 Web server listening at http://127.0.0.1:8765
    2015/05/02 13:32:49 Press ^C to stop

