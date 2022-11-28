
package transport

type Driver interface
{
    Send(namespace string, id uint32, message interface{}) error

    Commit() error
}
