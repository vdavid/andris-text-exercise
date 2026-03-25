import { deleteMessage } from '@/app/actions'
import type { Message } from '@/lib/supabase'

function formatTimestamp(isoString: string): string {
    const date = new Date(isoString)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}`
}

export default function MessageList({ messages }: { messages: Message[] }) {
    if (messages.length === 0) {
        return <p className="py-8 text-center text-sm text-gray-400">Még nincsenek üzenetek.</p>
    }

    return (
        <ul className="flex flex-col gap-3">
            {messages.map((message) => (
                <li key={message.id} className="flex items-start justify-between gap-3 rounded-lg border border-gray-200 p-4">
                    <div className="min-w-0 flex-1">
                        <p className="whitespace-pre-wrap break-words text-sm text-gray-800">{message.text}</p>
                        <time className="mt-1 block text-xs text-gray-400">{formatTimestamp(message.created_at)}</time>
                    </div>
                    <form action={deleteMessage}>
                        <input type="hidden" name="id" value={message.id} />
                        <button
                            type="submit"
                            className="shrink-0 rounded px-2 py-1 text-xs text-gray-400 transition-colors hover:bg-red-50 hover:text-red-600"
                            title="Törlés"
                        >
                            Törlés
                        </button>
                    </form>
                </li>
            ))}
        </ul>
    )
}
