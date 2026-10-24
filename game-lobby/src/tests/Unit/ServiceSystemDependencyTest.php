<?php

namespace Tests\Unit;

use PHPUnit\Framework\Attributes\DataProvider;
use PHPUnit\Framework\TestCase;
use ReflectionClass;
use ReflectionNamedType;

class ServiceSystemDependencyTest extends TestCase
{
    /** @return iterable<string, array{class-string}> */
    public static function gameServiceProvider(): iterable
    {
        yield from self::classesIn(dirname(__DIR__, 2).'/app/Services/Games', 'App\\Services\\Games', 'Service.php');
    }

    /** @return iterable<string, array{class-string}> */
    public static function gameSystemProvider(): iterable
    {
        yield from self::classesIn(dirname(__DIR__, 2).'/app/Systems/Games', 'App\\Systems\\Games', 'System.php');
    }

    #[DataProvider('gameServiceProvider')]
    public function test_game_services_do_not_depend_on_other_services(string $class): void
    {
        foreach ($this->constructorDependencyNames($class) as $dependency) {
            $this->assertFalse(
                str_ends_with($dependency, 'Service'),
                "{$class} must orchestrate Systems instead of another Service ({$dependency}).",
            );
        }
    }

    #[DataProvider('gameSystemProvider')]
    public function test_game_systems_do_not_depend_on_services_or_other_systems(string $class): void
    {
        foreach ($this->constructorDependencyNames($class) as $dependency) {
            $this->assertFalse(
                str_ends_with($dependency, 'Service') || str_ends_with($dependency, 'System'),
                "{$class} must be a self-contained functional component ({$dependency}).",
            );
        }
    }

    /** @return list<class-string> */
    private function constructorDependencyNames(string $class): array
    {
        $constructor = (new ReflectionClass($class))->getConstructor();
        if ($constructor === null) {
            return [];
        }

        $dependencies = [];
        foreach ($constructor->getParameters() as $parameter) {
            $type = $parameter->getType();
            if ($type instanceof ReflectionNamedType && ! $type->isBuiltin()) {
                /** @var class-string $dependency */
                $dependency = $type->getName();
                $dependencies[] = $dependency;
            }
        }

        return $dependencies;
    }

    /** @return iterable<string, array{class-string}> */
    private static function classesIn(string $directory, string $namespace, string $suffix): iterable
    {
        foreach (glob($directory.'/*'.$suffix) ?: [] as $file) {
            $class = $namespace.'\\'.pathinfo($file, PATHINFO_FILENAME);
            /** @var class-string $class */
            yield $class => [$class];
        }
    }
}
