# Journey
A blog engine written in Go, compatible with Ghost themes.

![Editor](https://raw.githubusercontent.com/kabukky/journey/gh-pages/images/journey.png)

## About
Please note that Journey is still in alpha and has not been tested in production. Not all Ghost theme helpers have been implemented yet so there is bound to be trouble with some Ghost themes. Please open a [New Issue](https://github.com/kabukky/journey/issues) if you encounter a theme that doesn't work.

#### Easy to work with
Create or update your posts from any place and any device. Simply point your browser to yourblog.url/admin/, log in, and start typing away!

#### Extensible
Write plugins in Lua to implement custom behavior when generating pages. Learn how to do it on the [Wiki](https://github.com/kabukky/journey/wiki/Creating-a-Journey-Plugin)!

#### Good stuff available right away
Use Ghost themes to design your blog. There's already a great community of designers working on Ghost compatible themes. Check out the [Ghost Marketplace](http://marketplace.ghost.org) to get an idea. You can also migrate your existing Ghost installation to Journey with [a few simple steps](https://github.com/kabukky/journey/wiki/Migrating-from-Ghost-to-Journey).

#### Good stuff to come
Hopefully. Planning the future of Journey, high priority goals are support of MySQL, PostgreSQL, and Google App Engine.

#### Easily secure
Other blog engines require you to install Nginx or Apache just to enable HTTPS. With Journey, simply enable HTTPS in the configuration and start using it for development purposes. For production, simply replace the generated certificates with your own and you are ready to go.

#### No dependencies
Don't worry about installing the correct version of Node.js, Python, or anything else. Just download the [latest release](https://www.github.com/kabukky/journey/releases) for your operating system and cpu architecture, then place the folder anywhere you like and run the Journey executable. Done!

#### Lightweight and fast
Journey is still in an early stage of development. However, initial tests indicate that it is much faster at generating pages than Ghost running on Node.js. It also eats very little of your precious memory. For example: Testing it on Mac OS X, it takes about 3.5 MB of it and then happily carries on doing its job.

This slimness makes Journey an ideal candidate for setting up micro blogs or hosting it on low-end vps machines or micro computers such as the Raspberry Pi.

#### Deployable anywhere
[Download the release package](https://www.github.com/kabukky/journey/releases) for Linux (AMD64, i386, ARM), Mac OS X (AMD64, i386) or Windows (AMD64, i386) and start using Journey right away. Build Journey from source to make it work on a multitude of other operating systems!

## Installing Journey
To get started with Journey, go to the the [Releases Page](https://github.com/kabukky/journey/releases) and download the zip file corresponding to your operating system and cpu architecture. Then extract Journey anywhere you like. Why not place it in your home folder (e.g. /home/youruser/journey/)?

After that, head over to [Setting up Journey](https://github.com/kabukky/journey/wiki/Setting-up-Journey) to configure your Journey blog on your local machine.

If you'd like to set up Journey on a Linux server, head over to [Installing Journey on Ubuntu Server](https://github.com/kabukky/journey/wiki/Installing-Journey-on-Ubuntu-Server) for a step-by-step tutorial.

Journey even runs as a Windows Azure Web App. It's a great way to try out or host a low traffic Journey blog for free on the internet! Head over to [Hosting Journey as a Windows Azure Web App](https://github.com/kabukky/journey/wiki/Hosting-Journey-as-a-Windows-Azure-Web-App) for a step-by-step tutorial.

## Plugins
Did you create a Journey plugin? Write me [@kabukky](https://twitter.com/kabukky) or me@kaihag.com and I'll add a link to it here.

## Questions?
Please read the [FAQ](https://github.com/kabukky/journey/wiki/FAQ) Wiki page or write to me@kaihag.com.

## Troubleshooting
Please refer to the [FAQ](https://github.com/kabukky/journey/wiki/FAQ) Wiki page if you experience any trouble running Journey.

If your issue isn't discussed there, please create a [New Issue](https://github.com/kabukky/journey/issues).

## Building from source
Please refer to the [Building Journey from source](https://github.com/kabukky/journey/wiki/Building-Journey-from-source) Wiki page for instructions on how to build Journey from source.

If you'd like to turn off the plugin system, you can use the build tag 'noplugins' to do so.

## Contributing to Journey
Pull requests are very much welcome. But please create them on the development branch. The master branch will only be updated for a new release.