// from http://www.typescriptlang.org/docs/tutorial.html
interface Person {
    firstName: string
    lastName: string
}

function greeter(person: Person): string {
    return `Hello, ${person.firstName} ${person.lastName}`;
}

const user = "Jane User";

console.log(greeter({ firstName: "Jon", lastName: "Doe" }));

//========================

const list: number[] = [1,2,3];

let people: Person[] = [];
people.push({ firstName: "Jane", lastName: "Doe" });
let x: [ string, number, boolean ] = ["a", 1, true];

enum Color { Red, Green, Blue };
const bgColor = Color.Green;


let whoKnows: any;
whoKnows = true;
whoKnows = 1;
whoKnows = "ok";

let strLength: number = (whoKnows as string).length;

//========================

function sum(...nums: number[]): number {
    return nums.reduce(((prev, n) => prev + n), 0);
}
console.log(sum(1,2,3));

//========================

function identity<T>(arg: T): T {
    return arg;
}
