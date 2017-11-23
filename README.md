# Moreau
Static website generator + deployer

**Yet another static website generator (SWG)?**

While this might be the first reaction, I would like to argue that I am trying to create something more.
Indeed there are many SWGs but how many of then can actually deploy the said created sites?

I am trying to take the concept of SWG to in a different direction by including deployment as part of the experience.

Let's see how it goes. So far current implementation is very limited - It works with GCP Storage
(because the official docs of AWS Static Website Hosting don't include any mention of using boto or aws sdk.).
It is good at blindly deploying without ensuring if it will overwrite existing essays.

I have tried to start this project multiple times and failed to make progress.
Sometimes creating crude minimum "Hello World App" does the trick :)

Hopefully **Moreau** will evolve over time to include things like

* More elaborate storage system.
* Support for AWS & SSH deployment
* Better UI
* More extensible html templating.

The points are included in descending order of importance.

However, hope this will become something useful.

~ Ent
