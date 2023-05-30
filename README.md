# Ad Server

In order to run this, download the project and run the following command in a terminal at project root:

`go build main.go && ./main`

The server will launch at localhost:8080

## High-Level Design

The project consists of three main components:

1. AdServer
2. Campaign Service
3. Router

### AdServer

AdServer receives Campaigns and is responsible for organizing them in way that will allow for fast and accurate
recommendations. It accomplishes this by using a sorted doubly-linked list. All active campaigns are inserted 
into this linked list. Each node, however, has several "next" and "previous" pointers which correspond to the
keywords they are associated with. At an abstract level, it appears like there are many linked lists within the
AdServer. The primary benefits of this setup are:

1. No duplication of data.
2. Removal of campaign from one list will remove it from all lists.
3. Adding campaign to one list will add it to all lists.
4. Constant look-ups for each keyword.

Campaigns are added / removed during their activation / expiration date using a regularly running async process.
When a new campaign is added, it is either immediately inserted into the underlying linkedlist or scheduled for 
insertion at its activation time. its removal is also scheduled this way. Each second the updater process will 
check for functions registered to that timestamp and execute them.

### Campaign Service

Campaign Service handles the definition, creation, and storage of Campaigns. In a production system this would
likely connect to a database, but in this implementation it just stores instances in memory.

### Router

Router is where all framework code lives and where interaction between AdServer and Campaign Service is coordinated.

## A note on the state of the project

I ended up becoming much busier than I had anticipated since we last spoke and unfortunately didn't get to finish
cleaning things up or writing some of the tests I would've liked to have written. Some of the inconsistencies are
a result of me learning as I was going and not having a chance to fix the old things to have a consisten style. 
These are all things we can discuss in detail over a video call though.
