import hello2 from "./hello-2";
import hello from "./hello";

/**
 * The idea for this file is to build a function that creates below output automagicly.
 * The function shoudl scan all dubdirectories and then add import the function and add
 * {"dirname": imported function} to below so that all functions can be exported automagicly.
 */

const functions = {
  hello: hello,
  "hello-2": hello2,
};

export default functions;
