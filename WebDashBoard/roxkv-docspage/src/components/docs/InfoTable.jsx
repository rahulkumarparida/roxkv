export default function InfoTable({ headers, rows, caption }) {
  return (
    <div className="my-6">
      {caption && <p className="text-sm text-[#9898ab] mb-2">{caption}</p>}
      <div className="overflow-x-auto rounded-xl border border-[#1e1e2e] bg-[#16161f]">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-[#1e1e2e] bg-[rgba(124,58,237,0.05)]">
              {headers.map((h, i) => (
                <th key={i} className="py-3 px-4 text-sm font-semibold text-[#a78bfa] whitespace-nowrap">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-[#1e1e2e]">
            {rows.map((row, i) => (
              <tr key={i} className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
                {row.map((cell, j) => (
                  <td key={j} className={`py-3 px-4 text-sm ${j === 0 ? 'font-mono text-white whitespace-nowrap' : 'text-[#9898ab]'}`}>
                    {cell}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
