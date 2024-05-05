## Initial list of tests
[None implemented yet :-() ]
* If non-admin user tries to execute admin-specific endpoints, flag error
# TEST
* Delete in frontend also deletes in backend
* Reset
  * Should only be accessible if admin is logged in
  * Should be possible to do it from admin dashboard
* From frontend, EITHER should forcibly refresh (if done by admin from frontend)
* OR should be impossible to access the old copy

## CHECK
* Refresh/populate the data of the logged in user (using that user's display list)  
* Refresh list of simulations when user logs in

### Authentication 
* Graceful failure if remote server not available DONE
* What does the refresh browser button do? 
