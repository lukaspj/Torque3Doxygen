if ($Game::argc != 2) {
    error("Should only have one argument");
    quit();
}

echo(eval($Game::argv[1]));

quit();