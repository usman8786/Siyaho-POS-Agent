export class AgentNotRunningError extends Error {
  constructor(message = 'Siyaho Printer Agent is not running.') {
    super(message);
    this.name = 'AgentNotRunningError';
  }
}

export class PrintAgentError extends Error {
  constructor(message, status) {
    super(message);
    this.name = 'PrintAgentError';
    this.status = status;
  }
}
