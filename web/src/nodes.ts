import { ClassicPreset, type GetSchemes } from 'rete'

/**
 * A single data socket for simplicity. You can split this
 * into stricter socket types (Domain vs List) and add connection
 * validation via editor pipes later if you want.
 */
export const dataSocket = new ClassicPreset.Socket('data')

/** Join: merges two string lists (deduped) */
export class JoinNode extends ClassicPreset.Node {
  constructor() {
    super('Join')
    this.addInput('a', new ClassicPreset.Input(dataSocket, 'List A'))
    this.addInput('b', new ClassicPreset.Input(dataSocket, 'List B'))
    this.addOutput('list', new ClassicPreset.Output(dataSocket, 'Merged'))
  }

  // inputs.a => [listA], inputs.b => [listB]
  data(inputs: { a?: string[][]; b?: string[][] }): { list: string[] } {
    const a = (inputs.a?.[0] ?? []) as string[]
    const b = (inputs.b?.[0] ?? []) as string[]
    const merged = Array.from(new Set([...a, ...b]))
    return { list: merged }
  }
}

/** Find Subdomains: Domain -> [subdomains...] (demo stub) */
export class FindSubdomainsNode extends ClassicPreset.Node {
  private domainText = 'example.com'

  constructor() {
    super('Find Subdomains')
    this.addInput('domain', new ClassicPreset.Input(dataSocket, 'Domain'))
    this.addOutput('list', new ClassicPreset.Output(dataSocket, 'Subdomains'))
    this.addControl(
      'domain',
      new ClassicPreset.InputControl('text', {
        initial: this.domainText,
        change: (v) => (this.domainText = String(v ?? ''))
      })
    )
  }

  data(inputs: { domain?: unknown[] }): { list: string[] } {
    const incoming = inputs.domain?.[0]
    const fallback = this.domainText.trim()
    const root = (typeof incoming === 'string' ? incoming : fallback).replace(/^https?:\/\//, '')
    if (!root) return { list: [] }
    // Deterministic demo data so you can wire graphs
    const out = [`www.${root}`, `api.${root}`, `dev.${root}`, `staging.${root}`]
    return { list: out }
  }
}

/** Find Wildcard Domains: Domain -> [wildcards...] (demo stub) */
export class FindWildcardDomainsNode extends ClassicPreset.Node {
  private domainText = 'example.com'

  constructor() {
    super('Find Wildcard Domains')
    this.addInput('domain', new ClassicPreset.Input(dataSocket, 'Domain'))
    this.addOutput('list', new ClassicPreset.Output(dataSocket, 'Wildcards'))
    this.addControl(
      'domain',
      new ClassicPreset.InputControl('text', {
        initial: this.domainText,
        change: (v) => (this.domainText = String(v ?? ''))
      })
    )
  }

  data(inputs: { domain?: unknown[] }): { list: string[] } {
    const incoming = inputs.domain?.[0]
    const fallback = this.domainText.trim()
    const root = (typeof incoming === 'string' ? incoming : fallback).replace(/^https?:\/\//, '')
    if (!root) return { list: [] }
    // Deterministic demo patterns
    const out = [`*.${root}`, `*.dev.${root}`, `*.staging.${root}`]
    return { list: out }
  }
}

export type Node = JoinNode | FindSubdomainsNode | FindWildcardDomainsNode
export type Conn = ClassicPreset.Connection<Node, Node>
export type Schemes = GetSchemes<Node, Conn>

/** Factory used by the palette */
export function createNode(kind: 'join' | 'find-subdomains' | 'find-wildcard-domains'): Node {
  switch (kind) {
    case 'join':
      return new JoinNode()
    case 'find-subdomains':
      return new FindSubdomainsNode()
    case 'find-wildcard-domains':
      return new FindWildcardDomainsNode()
    default:
      throw new Error(`Unknown node kind: ${kind}`)
  }
}
