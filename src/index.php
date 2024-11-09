<?php

function redirect($path)
{
	header('Location: ' . str_replace('/var/www/slcc/onlinedocs/', '', $path), true, 302);
	exit();
}

function page_content()
{
  $cmd = "cd " . getcwd() . "/onlinedocs/texi && make";
  shell_exec($cmd);

	$file = $_SERVER['REQUEST_URI'];

	if ($file == '/' || $file == '/master/' || $file == '/master') {
		$file = '/onlinedocs/master/index.html';
		redirect('/onlinedocs/master/index.html');
	}

	$path = getcwd() . '/onlinedocs/' . $file;

	if (!strrpos($path, '.html')) {
    $path = $path . '.html';
  }

  if (!file_exists($path)) {
    $path = getcwd() . '/onlinedocs/static/404.html';
  }
	
	echo file_get_contents($path);
}

function init()
{
	require 'src/template.php';
}
