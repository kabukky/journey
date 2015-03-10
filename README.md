# Journey
A blog engine written in Go, compatible with Ghost themes.

## About
Please note that Journey is still in alpha and has not been tested in production.

### Easy to work with
Create or update your posts from any place and any device. Simply point your browser to yourblogurl/admin, log in and start typing away!

### Good stuff to use right away
Use Ghost themes to design your blog. There's already a great community of designers working on Ghost compatible themes. Check out the [Ghost Marketplace](http://marketplace.ghost.org) to get an idea.

### Good stuff to come
Hopefully. Planning the future of Journey, high priority goals are: Plug-in support, MySQL and PostgreSQL support, and Google App Engine support.

### Easily secure
Other blog engines require you to install Nginx or Apache just to enable Https. With Journey, simply enable Https in the configuration and start using it for development purposes. For production, simply replace the generated certificates with your own and you are ready to go.

### No dependencies
Don't worry about installing the correct version of Node.js, Python, or anything else. Just download the [latest release](https://www.github.com/kabukky/journey/releases) for your operating system and cpu architecture, then place the folder anywhere you like and run the Journey executable. Done!

### Lightweight and fast
Journey is still in an early stage of development. However, initial tests indicate that it is about 10 times faster at generating pages than Ghost running on Node.js. It eats very little of your precious memory. For example: Testing it on a MacBook, it takes about 3.5 MB of it and then happily carries on doing its job.

This slimness makes Journey an ideal candidate for setting up micro blogs or hosting it on low-end vps machines or micro computers such as the Raspberry Pi.

### Deployable anywhere
[Download the release package](https://www.github.com/kabukky/journey/releases) for Linux (AMD64, i386, ARM), Mac OS X (AMD64, i386) or Windows (AMD64, i386) and start using Journey right away. Build Journey from source to make it work on a multitude of other operating systems!

## Installing Journey
Go to the the [Releases Page](https://github.com/kabukky/journey/releases) and download the zip file corresponding to your operating system and cpu architecture.

Then extract that zip file anywhere you like. You may also rename the extracted folder into "journey" if you so desire.

In the following section we will assume your Journey executable is located in /home/your-user/journey/

## Using Journey
### 1. Start Journey
In your Terminal, navigate to your Journey folder (e.g. /home/your-user/journey/) and start Journey by typing

    ./journey

when using Linux, Mac OS X, or another Unix

or

    journey.exe

when using Windows. Alternatively, you can just double-click on journey.exe

Journey is now running.

To visit your Journey blog, open

    http://127.0.0.1:8081

in your browser. You probably haven't written any blog posts to display yet. Let's change that.

Open

    http://127.0.0.1:8081/admin

in your browser. Fill out the information to create your Journey admin account. In the next step, log in using the user name and password you just provided.

From the admin area you can:
- create, edit, and delete blog posts
- edit your blog settings
- edit your user settings

### 2. Configure Journey
By editing the "config.json" file in your Journey root directory, you'll change the following settings:

**"HttpHostAndPort"**

This will change the port the Journey *http* server is listening on. If you don't want to bind to a particular ip address, writing just ":port number" as the value is fine (e.g. ":8081")

NOTE: If you want to change the port to 80 (HTTP default) you will probably have to set your firewall to redirect from port 80 to your Journey port or run Journey as root.

**"HttpsHostAndPort"**

This will change the port the Journey *https* server is listening on. If you don't want to bind to a particular ip address, writing just ":port number" as the value is fine (e.g. ":8082")

NOTE: If you want to change the port to 443 (HTTPS default) you will probably have to set your firewall to redirect from port 443 to your Journey port or run Journey as root.

**"HttpsUsage"**

This will change the https setting of your Journey blog. There are three possible values:
- "None"
  - Your Journey blog and admin area will only be accessible by http, NO https support is available.
- "AdminOnly"
  - Your admin area will ONLY be accessible by https (http connections will be redirected to https), your Journey blog will be accessible by both http and https.
- "All"
  - Your Journey blog and admin area will ONLY be accessible by https (http connections will be redirected to https).

NOTE: For a minimum of security, "HttpsUsage" should always be set to at least "AdminOnly" to ensure your login credentials and cookies are being sent using an encrypted connection.

When https is enabled ("HttpsUsage" is not set to "None"), you have to provide files containing a certificate and matching private key for the server.

Those files have to be placed in the content/https folder (e.g. /home/your-user/journey/content/https) as cert.pem and key.pem.

If Https is enabled and the cert.pem and cert.key files are not present in this directory, the application will generate new cert.pem and cert.key files upon startup.

Replace those files with your own as soon as possible and don't use them in production.

**"Url"**

This will change the url of your Journey blog. You have to change this to the host name the blog is supposed to be reachable under.

The "Url" setting is used to generate links to the blog (rss feeds and @blog.url helper) and to redirect incoming http connections to https.

### 3. Choose a theme
The Promenade theme is included by default to make Journey work out of the box. However, it is only intended to be used on a one author, personal website.

For a fully fledged, multiple author blog experience try the [Casper](https://github.com/TryGhost/Casper) theme from the makers of Ghost.

[Download it](https://github.com/TryGhost/Casper/releases) and place the Casper directory in content/themes folder (e.g. /home/your-user/journey/content/themes). Then select Casper from your admin panel under Settings/Blog.

After that, try some other themes! There's a whole world of Ghost themes out there. Find the one you like best.

### 4. Write your own theme
Finally, you can always write your own theme and use it with Journey. Start by visiting [http://themes.ghost.org](http://themes.ghost.org) and by reading one of the many tutorials that show you how to create a Ghost theme!

## Troubleshooting
### 1. "/lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.14' not found (required by journey)" when trying to start Journey
Your Linux distribution ships with an older version of libc6. Try updating libc6.

This may be tricky or not possible on Debian Wheezy stable. You can always try to compile Journey from source to link against your own version of libc6 (see below).

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
