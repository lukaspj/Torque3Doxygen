if ($Game::argc != 2) {
    error("Should have exactly one argument, received " @ $Game::argc);
    quit();
}

echo(eval($Game::argv[1]));

quit();