import { supabase } from '@/lib/supabase'
import type { Message } from '@/lib/supabase'
import MessageForm from '@/app/components/MessageForm'
import MessageList from '@/app/components/MessageList'

export default async function Home() {
    const { data: messages } = await supabase
        .from('messages')
        .select('*')
        .order('created_at', { ascending: false })

    return (
        <div className="mx-auto w-full max-w-[600px] px-4 py-10">
            <h1 className="mb-8 text-2xl font-bold text-gray-900">Üzenőfal</h1>
            <div className="mb-8">
                <MessageForm />
            </div>
            <MessageList messages={(messages as Message[]) ?? []} />
        </div>
    )
}
