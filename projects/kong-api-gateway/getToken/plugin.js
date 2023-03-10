'use strict';

const { execSync } = require('child_process');
 
class Custom {
  constructor(serverless, options) {
    this.serverless = serverless;
    this.options = options;
    this.commands = {};
    this.hooks = {};

    this.serverless.setProvider('Custom', this);

    this.stdin = process.stdin;
    this.stdout = process.stdout;
    this.stderr = process.stderr;

    this.defineCommands();
  }

  getEnv() {
     const service = this.serverless.service;
     return service.provider && service.provider.environment;
  }

  getConfig() {
    const service = this.serverless.service;
    return service.custom && service.custom.scripts;
  }

  defineCommands() {
    const config = this.getConfig();
    const commands = config && config.commands;
    if (!commands) return;

    for (const name of Object.keys(commands)) {
      if (!this.commands[name]) {
        this.commands[name] = { lifecycleEvents: [] };
      }
      this.commands[name].lifecycleEvents.push(name);

      this.hooks[`${name}:${name}`] = this.runCommand.bind(this, name);
    }
  }
  
    runCommand(name) {
    const commands = this.getConfig().commands;
    const command = commands[name];
    this.execute(command);
  }

  execute(command) {
    execSync(command, { env: { ...process.env, ...this.getEnv() }, stdio: [this.stdin, this.stdout, this.stderr] })
  }
}  

module.exports = Custom;
