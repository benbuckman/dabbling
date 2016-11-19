import sys


def understand_errors():

    try:
        n = 10 / 0
        print("After failure, should not appear!")
    except ZeroDivisionError as err:
        print('Failed to divide: %s' % err)
    except:
        print("Unexpected error:", sys.exc_info()[0])

    n = 10 / 1

    print(int(n))


understand_errors()
