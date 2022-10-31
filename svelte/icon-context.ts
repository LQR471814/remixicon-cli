export const iconKey = Symbol("web-stdlib")
export type Context = Partial<{
  className: string
  width: number
  height: number
  x: number
  y: number
  type: "line" | "fill"
}>
