# Unsystematic list of required changes and fixes.  
## Semi-immediate
* At startup does not connect to server (This is because of Heroku ephemeral file system so we have to switch to PostGres unfortunately)
* We don't know when and why it switches to the player selection at first use
* Once in as user guest, and switch to another player, we are returned to the dashboard as user guest.
* When we then refresh we get an error that we cannot play as whoever we select
* But we can then go to the user dashboard and it's OK

## Some time
* Get 'static' working  
* Need a more efficient method for accessing links between objects (eg maps, or a local database)
* Switch gin to production mode
* Parameterise the templates to reduce boilerplate
  
## Admin
* complete admin dashboard  
* Should be able to delete users 
* Should be able to delete user simulations 
* Should be able to view user simulations

