# Journey
A blog engine written in Go, compatible with Ghost themes.

## Installing Journey
Go to the the [Releases Page](https://github.com/kabukky/journey/releases) and download the zip file corresponding to your operating system and cpu architecture.

Then extract that zip file anywhere you like. You may also rename the extracted folder into "journey" if you like.

In the following section we will assume your Journey executable is located in /home/your-user/journey/

## Using Journey

### 1. First start
In your Terminal, navigate to your Journey folder (e.g. /home/your-user/journey/) and start Journey by typing

    ./journey

when using Linux, Mac OS X or another Unix

or

    journey.exe

when using Windows.

Journey is now running.

To visit your Journey blog, open

    http://127.0.0.1:8081

in your browser. You probably haven't written any blog posts to display yet. Let's change that.

Open

    http://127.0.0.1:8081/admin

in your browser. Fill out the information to create your Journey admin account. In the next step, log in using the user name and password you just provided.

In the admin area you can:
- create, edit, and delete blog posts
- edit your blog settings
- edit your user settings

### 2. Configuration

By editing the "config.json" file you can change the following settings:

**"HttpHostAndPort"**

This will change the port the Journey *http* server is listening on. If you don't want to bind to a particular ip address, writing just ":port number" as the value is fine (e.g. ":80")

NOTE: If you change the port to 80 (HTTP default) you will probably have to run Journey as root.

**"HttpsHostAndPort"**

This will change the port the Journey *https* server is listening on. If you don't want to bind to a particular ip address, writing just ":port number" as the value is fine (e.g. ":443")

NOTE: If you change the port to 443 (HTTPS default) you will probably have to run Journey as root.

**"HttpsUsage"**

This will change the https setting of your Journey blog. There are three possible values:
- "None"
  - Your Journey blog and admin area will only be accessible by http, NO https support is available.
- "AdminOnly"
  - Your admin area will ONLY be accessible by https (http connections will be redirected to https), your Journey blog will be accessible by both http and https.
- "All"
  - Your Journey blog and admin area will ONLY be accessible by https (http connections will be redirected to https).

NOTE: For a minimum of security, "HttpsUsage" should always be set to at least "AdminOnly" to ensure your login credentials and cookies are being sent using an encrypted connection.

**"Url"**

This will change the url of your Journey blog. You have to change this to the host name the blog is supposed to be reachable under.

The "Url" setting will be used to generate links to the blog (rss feeds and @blog.url helper) and to redirect incoming http connections to https.

### 3. Choose a theme

The Promenade theme is included by default to make Journey work out of the box. However, it is only intended to be used on a one author, personal website.

For a fully fledged, multiple author blog experience try the [Casper](https://github.com/TryGhost/Casper) theme from the makers of Ghost.

[Download it](https://github.com/TryGhost/Casper/releases) and place the Casper directory in content/themes folder (e.g. /home/your-user/journey/content/themes). Then select Casper from your admin panel under Settings/Blog.

Then try some other themes! There's a whole world of Ghost themes out there. Find the one you like best.

### 3. Write your own theme

Finally, you can always write your own theme and use it with Journey. Start by visiting [http://themes.ghost.org]{http://themes.ghost.org) and by reading one of the many tutorials that show you how to create a Ghost theme!

## Building from source
Prerequisites
- Install Go if it's not on your system already and set the correct GOPATH.
- Install Git if it's not on your system already.

Then run

    go get github.com/kabukky/journey
  
This will download the source into $GOPATH/src/github.com/kabukky/journey.

In your terminal, change into that journey folder and run

    git submodule update --init --recursive

This will download the default theme into $GOPATH/src/github.com/kabukky/journey/content/themes.

Still in that journey folder run

    go build

This will build the Journey binary.

You may copy the "journey" binary file, the "content" folder, the "built-in" folder and the "config.json" file to a new location (e. g. /home/your-user/journey/). Then run the Journey binary from that new location to start the server.
