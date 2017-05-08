function greeter(person) {
    return `Hello, ${person.firstName} ${person.lastName}`;
}
const user = "Jane User";
console.log(greeter({ firstName: "Jon", lastName: "Doe" }));
//========================
const list = [1, 2, 3];
let people = [];
people.push({ firstName: "Jane", lastName: "Doe" });
let x = ["a", 1, true];
var Color;
(function (Color) {
    Color[Color["Red"] = 0] = "Red";
    Color[Color["Green"] = 1] = "Green";
    Color[Color["Blue"] = 2] = "Blue";
})(Color || (Color = {}));
;
const bgColor = Color.Green;
let whoKnows;
whoKnows = true;
whoKnows = 1;
whoKnows = "ok";
let strLength = whoKnows.length;
//========================
function sum(...nums) {
    return nums.reduce(((prev, n) => prev + n), 0);
}
console.log(sum(1, 2, 3));
//========================
function identity(arg) {
    return arg;
}
//# sourceMappingURL=index.js.map