# Recipe App in Go 
Frontend: \
/             -to get Home page \
/account      -to get Account page \
/recipe/:id   -to get page of a Recipe \
/tutorials    -to get the page of a tutorial (Currently it always shows the same video)\

/get/home     -to get content of Home page \
/get/account  -to get content of Account page \

Routes: \
GET    /get          -to get every Recipe \
POST   /create       -to create a Recipe \
GET    /getbyid/:id  -to get a Recipe by id \
DELETE /delete/:id   -to delete a Recipe by id \
GET    /select/:id   -to select a Recipe \
GET    /deselect/:id -to deselect a Recipe \
GET    /colormode    -to get/set the colormode(darkmode|lightmode) \
POST   /filter       -to filter for recipes and cookingtime

# Next Update
/recomendet   -better recomendation 