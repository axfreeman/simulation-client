# Frontend for the simulation of a capitalist economy  
* Interface with remote server via its API.  
* This version uses http://github.com/axfreeman/simulation-api-with-api-key-login    
* There is more documentation at http://github.com/axfreeman/simulation-api    
     
## Philosophy
The aim of this version is to focus on the UI logic by reducing the 
overhead of authentication to the bare minimum possible.  

The backend accepts only api-key header requests.  

This client should use an api-key issued by the developer and run as 
a standalone app.

To scale up, we may introduce a local login system which consults the 
server to ascertain which keys have been issued and asks the user 
to choose one.  

