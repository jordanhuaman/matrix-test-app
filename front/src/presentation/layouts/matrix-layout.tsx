import type { ReactNode } from "react"

export function MatrixLayout({ sidebar, children }: { sidebar: ReactNode; children: ReactNode }) {
  return (
    <div className="flex h-[calc(100vh-4rem)]">
      <aside className="w-80 border-r border-border bg-card overflow-y-auto shrink-0 hidden lg:block">
        {sidebar}
      </aside>
      <div className="flex-1 overflow-y-auto p-6">
        {children}
      </div>
    </div>
  )
}
