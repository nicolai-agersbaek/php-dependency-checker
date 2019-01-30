<?php
declare(strict_types = 1);

//printClasses();
//printConstants();

//echo 'php -r $F=get_defined_functions(true)["internal"];foreach ($F as $f) {printf("$f\n");}';

$eol = PHP_EOL;

$eolRepl = [
    "\n" => '\n',
    "\r" => '\r',
];
//$eolFormatted = \strtr(PHP_EOL, $eolRepl);


$eolFormatted=strtr(PHP_EOL,["\n"=>'\n',"\r"=>'\r']);

//$eolFormatted=str_replace(["\r","\n"],['\r','\n'],PHP_EOL);

echo "PHP_EOL = {$eolFormatted}\n";

function printFunctions() : void
{
    /** @noinspection PotentialMalwareInspection */
    $definedFunctions = get_defined_functions(true);
    $internalFunctions = $definedFunctions['internal'] ?? [];
    
    $batches = batch($internalFunctions, 10);
    
    foreach ($batches[0] as $function) {
        \printf("%s\n", $function);
    }
    
    printf("last x:\n");
    
    foreach (end($batches) as $function) {
        \printf("%s\n", $function);
    }
}

function printClasses() : void
{
    $classes = get_declared_classes();
    $interfaces = get_declared_interfaces();
    $traits = get_declared_traits();
    
    printLinesWithTitle('CLASSES:', $classes);
    printLinesWithTitle('INTERFACES:', $interfaces);
    printLinesWithTitle('TRAITS:', $traits);
}

function getDefinedConstants()
{
    $byCategory = get_defined_constants(true);
    
    // Remove user-defined constants and merge the built-in constants into a single array
    unset($byCategory['user']);
    
    $allConstants = [];
    
    foreach ($byCategory as $C) {
        $allConstants[] = array_keys($C);
    }
    
    return array_merge(...$allConstants);
}

function printConstants() : void
{
    printLinesWithTitle('CONSTANTS:', getDefinedConstants());
}

/**
 * @param array $items
 * @param int   $batchSize
 *
 * @return array
 */
function batch(array $items, int $batchSize) : array
{
    $numBatches = \ceil(\count($items) / $batchSize);
    $batches = [];
    
    /** @noinspection PhpVariableNamingConventionInspection */
    for ($i = 0; $i < $numBatches; $i++) {
        $batches[] = array_slice($items, $i * $batchSize, $batchSize);
    }
    
    return $batches;
}

/**
 * @param string ...$lines
 */
function printLines(string ...$lines) : void
{
    \printf(\implode("\n", $lines));
}

/**
 * @param string $title
 */
function printTitle(string $title) : void
{
    $separator = str_repeat('-', strlen($title));
    
    \printf("\n%s\n", $separator);
    \printf("%s\n", $title);
    \printf("%s\n", $separator);
}

/**
 * @param string $title
 * @param array  $lines
 */
function printLinesWithTitle(string $title, array $lines) : void
{
    printTitle($title);
    printLines(...$lines);
}
